<!-- /src/lib/components/ChatBubble.svelte -->
<script lang="ts">
    import { onMount, onDestroy } from 'svelte';

    let isOpen = false;
    let message = '';
    let messages: { role: 'user' | 'assistant' | 'reasoning'; content: string }[] = [
        { content: 'How can I help you today?', role: 'assistant' }
    ];
    let ws: WebSocket | null = null;
    let fullBuffer = ''; // Buffer to accumulate all incoming words

    // Simple markdown parser for basic formatting
    function formatMessage(text: string): string {
        text = text.replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>');
        text = text.replace(/(\d+\.\s+.*?(?=\n|$))/g, '<li>$1</li>');
        if (text.includes('<li>')) {
            text = `<ol>${text}</ol>`;
        }
        text = text.replace(/```(.*?)```/gs, '<pre><code>$1</code></pre>');
        text = text.replace(/\n/g, '<br>');
        return text;
    }

    // Process the full buffer to extract reasoning and response
    function processBuffer(buffer: string) {
        const thinkStart = buffer.indexOf('<think>');
        const thinkEnd = buffer.indexOf('</think>');

        // If we have a complete <think> section
        if (thinkStart !== -1 && thinkEnd !== -1 && thinkEnd > thinkStart) {
            const reasoningText = buffer.slice(thinkStart + 7, thinkEnd).trim();
            const responseText = buffer.slice(thinkEnd + 8).trim();

            // Check if the last messages are from this response cycle
            let lastWasReasoning = messages[messages.length - 1]?.role === 'reasoning';
            let lastWasAssistant = messages[messages.length - 1]?.role === 'assistant';
            let secondLastWasReasoning = messages[messages.length - 2]?.role === 'reasoning';

            // Update or append messages without clearing history
            if (reasoningText && responseText) {
                if (lastWasAssistant && secondLastWasReasoning) {
                    // Update both reasoning and response
                    messages = [
                        ...messages.slice(0, -2),
                        { role: 'reasoning', content: formatMessage(reasoningText) },
                        { role: 'assistant', content: formatMessage(responseText) }
                    ];
                } else if (lastWasReasoning) {
                    // Update reasoning and add response
                    messages = [
                        ...messages.slice(0, -1),
                        { role: 'reasoning', content: formatMessage(reasoningText) },
                        { role: 'assistant', content: formatMessage(responseText) }
                    ];
                } else {
                    // Append both
                    messages = [
                        ...messages,
                        { role: 'reasoning', content: formatMessage(reasoningText) },
                        { role: 'assistant', content: formatMessage(responseText) }
                    ];
                }
            } else if (reasoningText) {
                if (lastWasReasoning) {
                    messages = [
                        ...messages.slice(0, -1),
                        { role: 'reasoning', content: formatMessage(reasoningText) }
                    ];
                } else {
                    messages = [
                        ...messages,
                        { role: 'reasoning', content: formatMessage(reasoningText) }
                    ];
                }
            }
        } else if (thinkStart !== -1) {
            // Partial think section, show as reasoning
            const reasoningText = buffer.slice(thinkStart + 7).trim();
            if (messages[messages.length - 1]?.role === 'reasoning') {
                messages = [
                    ...messages.slice(0, -1),
                    { role: 'reasoning', content: formatMessage(reasoningText) }
                ];
            } else {
                messages = [
                    ...messages,
                    { role: 'reasoning', content: formatMessage(reasoningText) }
                ];
            }
        } else {
            // No think section yet, treat as response
            if (messages[messages.length - 1]?.role === 'assistant') {
                messages = [
                    ...messages.slice(0, -1),
                    { role: 'assistant', content: formatMessage(buffer.trim()) }
                ];
            } else {
                messages = [
                    ...messages,
                    { role: 'assistant', content: formatMessage(buffer.trim()) }
                ];
            }
        }
    }

    function connectWebSocket() {
        ws = new WebSocket('ws://127.0.0.1:9089/ws');

        ws.onopen = () => {
            console.log('WebSocket connected');
        };

        ws.onmessage = (event) => {
            const word = event.data;

            // Handle error messages immediately
            if (word.includes('Error')) {
                messages = [...messages, { role: 'assistant', content: word }];
                fullBuffer = '';
                return;
            }

            // Accumulate all words
            fullBuffer += (fullBuffer ? ' ' : '') + word;

            // Process the buffer to update messages
            processBuffer(fullBuffer);
        };

        ws.onerror = (error) => {
            console.error('WebSocket error:', error);
            messages = [
                ...messages,
                { role: 'assistant', content: 'Error connecting to assistant' }
            ];
            fullBuffer = '';
        };

        ws.onclose = () => {
            console.log('WebSocket closed');
            ws = null;
            fullBuffer = '';
        };
    }

    function sendMessage() {
        if (!ws || ws.readyState !== WebSocket.OPEN) {
            messages = [...messages, { role: 'assistant', content: 'Assistant is offline' }];
            return;
        }
        if (message.trim() === '') return;

        ws.send(message);
        messages = [...messages, { role: 'user', content: message }];
        fullBuffer = '';
        message = '';
    }

    function toggleChat() {
        isOpen = !isOpen;
        if (isOpen && !ws) {
            connectWebSocket();
        }
    }

    onDestroy(() => {
        if (ws) {
            ws.close();
        }
    });

    function handleKeydown(event: KeyboardEvent) {
        if (event.key === 'Enter' && !event.shiftKey) {
            event.preventDefault();
            sendMessage();
        }
    }
</script>

<div class="fixed bottom-6 right-6 z-50">
    <button
        aria-label="chat"
        on:click={toggleChat}
        class="bg-retroTeal text-retroCream w-14 h-14 rounded-full flex items-center justify-center shadow-lg hover:bg-retroCoral transition-all duration-300 border-2 border-retroBlack cursor-pointer"
        title="Chat with Assistant"
    >
        <svg
            xmlns="http://www.w3.org/2000/svg"
            class="h-6 w-6"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
        >
            <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z"
            />
        </svg>
    </button>

    {#if isOpen}
        <div
            class="absolute bottom-16 right-0 w-[800px] h-[600px] bg-retroCream rounded-lg shadow-xl border-4 border-retroGray flex flex-col"
        >
            <div
                class="bg-retroGray text-retroCream p-2 rounded-t-md font-pixel text-sm flex justify-between items-center border-b-2 border-retroBlack"
            >
                <span>Store Assistant</span>
                <button
                    on:click={toggleChat}
                    class="text-retroCream hover:text-retroCoral transition-colors cursor-pointer"
                >
                    ✕
                </button>
            </div>

            <div class="flex-1 p-3 overflow-y-auto">
                {#each messages as msg}
                    <div class="{msg.role === 'user' ? 'text-right' : 'text-left'} mb-2">
                        <span
                            class="{msg.role === 'user'
                                ? 'bg-retroTeal text-retroCream'
                                : msg.role === 'reasoning'
                                  ? 'bg-gray-200 text-gray-700 italic'
                                  : 'bg-retroGray text-retroCream'} inline-block p-2 rounded-md font-body text-sm"
                        >
                            {#if msg.role === 'assistant' || msg.role === 'reasoning'}
                                {@html msg.content}
                            {:else}
                                {msg.content}
                            {/if}
                        </span>
                        {#if msg.role === 'reasoning'}
                            <small class="block text-gray-500 mt-1">Model Reasoning</small>
                        {/if}
                    </div>
                {/each}
            </div>

            <div class="p-3 border-t-2 border-retroBlack">
                <textarea
                    bind:value={message}
                    on:keydown={handleKeydown}
                    placeholder="Type a message..."
                    class="w-full h-12 p-2 bg-retroCream text-retroGray border-2 border-retroBlack rounded-md font-body resize-none focus:outline-none focus:border-retroTeal"
                ></textarea>
                <button
                    on:click={sendMessage}
                    class="mt-2 w-full bg-retroCoral text-retroCream py-1 rounded-md font-pixel text-sm hover:bg-retroTeal transition-all duration-300 border-2 border-retroBlack cursor-pointer"
                >
                    Send
                </button>
            </div>
        </div>
    {/if}
</div>

<style>
    strong {
        font-weight: bold;
    }
    ol {
        padding-left: 1.5em;
        margin: 0.5em 0;
    }
    li {
        margin: 0.2em 0;
    }
    pre {
        background: #eee;
        padding: 0.5em;
        border-radius: 4px;
        overflow-x: auto;
    }
    code {
        font-family: monospace;
    }
    .italic {
        font-style: italic;
    }
</style>
