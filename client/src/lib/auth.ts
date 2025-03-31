// /src/lib/auth.ts
import { writable } from 'svelte/store';

export const auth = writable<{
    token: string | null;
    user: {
        id: number;
        email: string;
        firstName: string;
        lastName: string;
        role: string;
        phone: string;
    } | null;
}>({
    token: null,
    user: null
});

if (typeof window !== 'undefined') {
    const storedToken = localStorage.getItem('token');
    const storedUser = localStorage.getItem('user');
    if (storedToken && storedUser) {
        auth.set({ token: storedToken, user: JSON.parse(storedUser) });
    }
}

auth.subscribe(({ token, user }) => {
    if (typeof window !== 'undefined') {
        if (token && user) {
            localStorage.setItem('token', token);
            localStorage.setItem('user', JSON.stringify(user));
        } else {
            localStorage.removeItem('token');
            localStorage.removeItem('user');
        }
    }
});

export function logout() {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    auth.set({ token: null, user: null });
}
