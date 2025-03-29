<!-- /src/routes/admin/products/+page.svelte -->
<script lang="ts">
    import Button from '$lib/components/Button.svelte';
    import { onMount } from 'svelte';
    import { createUploader } from '$lib/utils/uploadthing';
    import { UploadButton } from '@uploadthing/svelte';
    import type { OurFileRouter } from '$lib/server/uploadthing';
    import { GLOBAL } from '$lib';

    interface ProductImage {
        url: string;
        isPrimary: boolean;
        order: number;
    }

    interface ProductSpec {
        key: string;
        value: string;
    }

    interface Product {
        name: string;
        price: number;
        description: string;
        categoryId: string;
        stock: number;
        images: ProductImage[];
        specs: ProductSpec[];
    }

    interface Category {
        id: number;
        name: string;
        description: string;
    }

    let product: Product = {
        name: '',
        price: 0,
        description: '',
        categoryId: '',
        stock: 0,
        images: [{ url: '', isPrimary: true, order: 0 }],
        specs: [{ key: '', value: '' }]
    };

    let categories: Category[] = [];
    let error: string = '';
    let success: string = '';
    let uploadedImageUrl: string = '';
    let isSubmitting: boolean = false;

    const uploader = createUploader('imageUploader', {
        onClientUploadComplete: (res) => {
            if (res && res.length > 0) {
                const newImageUrl = res[0].url;
                uploadedImageUrl = newImageUrl;
                if (product.images[0].url === '') {
                    product.images[0].url = newImageUrl;
                } else {
                    product.images = [
                        ...product.images,
                        { url: newImageUrl, isPrimary: false, order: product.images.length }
                    ];
                }
                success = 'Image uploaded successfully!';
            }
        },
        onUploadError: (error: Error) => {
            error = Error(`Upload failed: ${error.message}`);
        },
        onUploadBegin: () => {
            success = 'Uploading image...';
        }
    });

    onMount(async () => {
        try {
            const response = await fetch(`${GLOBAL.SERVER_URL}/categories`);
            if (response.ok) {
                const res = await response.json();
                categories = res.categories;
            } else {
                error = 'Failed to load categories';
            }
        } catch (err) {
            error = 'An error occurred while fetching categories';
        }
    });

    function addImageField(): void {
        product.images = [
            ...product.images,
            { url: '', isPrimary: false, order: product.images.length }
        ];
    }

    function addSpecField(): void {
        product.specs = [...product.specs, { key: '', value: '' }];
    }

    function copyToClipboard(): void {
        navigator.clipboard.writeText(uploadedImageUrl);
        success = 'URL copied to clipboard!';
    }

    async function submitProduct(): Promise<void> {
        error = '';
        success = '';
        isSubmitting = true;

        // Client-side validation
        if (product.name === '') {
            error = 'Product name is required';
            isSubmitting = false;
            return;
        }
        if (product.price <= 0) {
            error = 'Price must be greater than 0';
            isSubmitting = false;
            return;
        }
        if (product.description === '') {
            error = 'Description is required';
            isSubmitting = false;
            return;
        }
        if (product.categoryId === '') {
            error = 'Category is required';
            isSubmitting = false;
            return;
        }
        if (product.stock < 0) {
            error = 'Stock cannot be negative';
            isSubmitting = false;
            return;
        }
        if (!product.images.some((img) => img.url)) {
            error = 'At least one image is required';
            isSubmitting = false;
            return;
        }
        if (!product.images.some((img) => img.isPrimary)) {
            error = 'At least one image must be marked as primary';
            isSubmitting = false;
            return;
        }

        try {
            const response = await fetch(`${GLOBAL.SERVER_URL}/products`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    name: product.name,
                    price: parseFloat(product.price.toString()),
                    description: product.description,
                    categoryId: parseInt(product.categoryId),
                    stock: parseInt(product.stock.toString()),
                    images: product.images,
                    specs: product.specs.filter((spec) => spec.key && spec.value) // Filter out empty specs
                })
            });

            if (response.ok) {
                const data = await response.json();
                success = data.message;
                product = {
                    name: '',
                    price: 0,
                    description: '',
                    categoryId: '',
                    stock: 0,
                    images: [{ url: '', isPrimary: true, order: 0 }],
                    specs: [{ key: '', value: '' }]
                };
                uploadedImageUrl = '';
            } else {
                const data = await response.json();
                error = data.error || 'Failed to add product';
            }
        } catch (err) {
            error = 'An error occurred while adding the product';
        } finally {
            isSubmitting = false;
        }
    }
</script>

<main class="container mx-auto px-4 py-8 bg-retroCream min-h-screen">
    <section class="max-w-3xl mx-auto">
        <h1 class="font-pixel text-3xl text-retroGray mb-6 text-center animate-pulse">
            Add New Retro Gear
        </h1>

        {#if error}
            <div class="bg-retroCoral text-retroCream p-4 rounded-lg mb-6 font-pixel">{error}</div>
        {/if}
        {#if success}
            <div class="bg-retroTeal text-retroCream p-4 rounded-lg mb-6 font-pixel">{success}</div>
        {/if}

        <!-- Image Upload Section -->
        <div
            class="bg-white p-6 rounded-lg border-4 border-retroBlack shadow-[4px_4px_0_#26A69A] mb-6"
        >
            <h2 class="font-pixel text-xl text-retroGray mb-4">Upload Image</h2>
            <UploadButton
                {uploader}
                class="ut-button bg-retroBlue text-retroCream hover:bg-retroCoral font-pixel px-3 py-1.5 rounded-md text-xs cursor-pointer"
            >
                <span slot="button-content" let:state>
                    {state.isUploading ? 'Uploading...' : 'Upload Image'}
                </span>
            </UploadButton>
            {#if uploadedImageUrl}
                <div class="mt-4">
                    <p class="font-pixel text-retroGray">Latest Uploaded Image URL:</p>
                    <input
                        type="text"
                        value={uploadedImageUrl}
                        readonly
                        class="w-full p-2 border-2 border-retroGray bg-retroCream font-pixel text-retroBlack rounded mb-2"
                    />
                    <Button
                        variant="secondary"
                        subClass="bg-retroPlum text-retroCream hover:bg-retroPurple"
                        onClick={copyToClipboard}
                    >
                        Copy URL
                    </Button>
                </div>
            {/if}
        </div>

        <form on:submit|preventDefault={submitProduct} class="space-y-6">
            <!-- Basic Product Info -->
            <div
                class="bg-white p-6 rounded-lg border-4 border-retroBlack shadow-[4px_4px_0_#26A69A]"
            >
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <label class="font-pixel text-retroGray block mb-2">Product Name</label>
                        <input
                            type="text"
                            bind:value={product.name}
                            required
                            class="w-full p-2 border-2 border-retroGray bg-retroCream font-pixel text-retroBlack rounded"
                        />
                    </div>
                    <div>
                        <label class="font-pixel text-retroGray block mb-2">Price ($)</label>
                        <input
                            type="number"
                            step="0.01"
                            bind:value={product.price}
                            required
                            class="w-full p-2 border-2 border-retroGray bg-retroCream font-pixel text-retroBlack rounded"
                        />
                    </div>
                </div>

                <div class="mt-4">
                    <label class="font-pixel text-retroGray block mb-2">Description</label>
                    <textarea
                        bind:value={product.description}
                        required
                        class="w-full p-2 border-2 border-retroGray bg-retroCream font-pixel text-retroBlack rounded h-32"
                    ></textarea>
                </div>

                <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mt-4">
                    <div>
                        <label class="font-pixel text-retroGray block mb-2">Category</label>
                        <select
                            bind:value={product.categoryId}
                            required
                            class="w-full p-2 border-2 border-retroGray bg-retroCream font-pixel text-retroBlack rounded"
                        >
                            <option value="">Select Category</option>
                            {#each categories as category}
                                <option value={category.id}>{category.name}</option>
                            {/each}
                        </select>
                    </div>
                    <div>
                        <label class="font-pixel text-retroGray block mb-2">Stock Quantity</label>
                        <input
                            type="number"
                            bind:value={product.stock}
                            required
                            class="w-full p-2 border-2 border-retroGray bg-retroCream font-pixel text-retroBlack rounded"
                        />
                    </div>
                </div>
            </div>

            <!-- Images -->
            <div
                class="bg-white p-6 rounded-lg border-4 border-retroBlack shadow-[4px_4px_0_#26A69A]"
            >
                <h2 class="font-pixel text-xl text-retroGray mb-4">Product Images</h2>
                <p class="font-pixel text-retroGray mb-4">Uploaded images appear here:</p>
                {#each product.images as image, index}
                    <div class="flex items-center gap-4 mb-4">
                        <input
                            type="text"
                            bind:value={image.url}
                            placeholder="Image URL"
                            required={index === 0}
                            class="flex-1 p-2 border-2 border-retroGray bg-retroCream font-pixel text-retroBlack rounded"
                        />
                        <label class="flex items-center gap-2 font-pixel text-retroGray">
                            <input
                                type="checkbox"
                                bind:checked={image.isPrimary}
                                on:change={() => {
                                    if (image.isPrimary) {
                                        product.images.forEach((img) => {
                                            if (img !== image) img.isPrimary = false;
                                        });
                                    }
                                }}
                            />
                            Primary
                        </label>
                    </div>
                {/each}
                <Button
                    variant="secondary"
                    subClass="bg-retroBlue text-retroCream hover:bg-retroCoral"
                    onClick={addImageField}
                >
                    Add Another Image
                </Button>
            </div>

            <!-- Specs -->
            <div
                class="bg-white p-6 rounded-lg border-4 border-retroBlack shadow-[4px_4px_0_#26A69A]"
            >
                <h2 class="font-pixel text-xl text-retroGray mb-4">Product Specs</h2>
                {#each product.specs as spec}
                    <div class="grid grid-cols-2 gap-4 mb-4">
                        <input
                            type="text"
                            bind:value={spec.key}
                            placeholder="Spec Key (e.g., Memory)"
                            class="p-2 border-2 border-retroGray bg-retroCream font-pixel text-retroBlack rounded"
                        />
                        <input
                            type="text"
                            bind:value={spec.value}
                            placeholder="Spec Value (e.g., 2MB)"
                            class="p-2 border-2 border-retroGray bg-retroCream font-pixel text-retroBlack rounded"
                        />
                    </div>
                {/each}
                <Button
                    variant="secondary"
                    subClass="bg-retroBlue text-retroCream hover:bg-retroCoral"
                    onClick={addSpecField}
                >
                    Add Another Spec
                </Button>
            </div>

            <!-- Submit -->
            <div class="text-center">
                <Button
                    variant="primary"
                    subClass="bg-retroCoral text-retroCream hover:bg-retroPurple hover:scale-105 transition-all"
                    type="submit"
                    disabled={isSubmitting}
                >
                    {isSubmitting ? 'Adding...' : 'Add Product to Inventory'}
                </Button>
            </div>
        </form>
    </section>
</main>

<style>
    .font-pixel {
        font-family: 'Jura', sans-serif;
    }
</style>
