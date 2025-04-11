// src/types/stripe.d.ts
interface Stripe {
    elements: () => any;
    confirmCardPayment: (clientSecret: string, data: any) => Promise<any>;
}

declare global {
    interface Window {
        Stripe: (publishableKey: string) => Stripe;
    }
}
