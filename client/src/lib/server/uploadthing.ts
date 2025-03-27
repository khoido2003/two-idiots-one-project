import { createUploadthing } from 'uploadthing/server';
import type { FileRouter } from 'uploadthing/server';

const f = createUploadthing();

const auth = (req: Request) => ({ id: 'fakeId' });

export const ourFileRouter = {
    imageUploader: f({
        image: {
            maxFileSize: '4MB',
            maxFileCount: 1
        }
    })
        .middleware(async ({ req }) => {
            // Add authentication later
            const user = await auth(req);
            if (!user) throw new Error('Unauthorized');
            return { userId: user.id };
        })
        .onUploadComplete(async ({ metadata, file }) => {
            console.log('Upload complete for userId:', metadata.userId);
            console.log('file url', file.url); // Note: `ufsUrl` might be `url` depending on version
        })
} satisfies FileRouter;

export type OurFileRouter = typeof ourFileRouter;
