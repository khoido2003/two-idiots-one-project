// api/uploadthing.js
import { createUploadthing } from "uploadthing/express";

const f = createUploadthing();

export const uploadRouter = {
  imageUploader: f({
    image: {
      maxFileSize: "4MB",
      maxFileCount: 1,
    },
  }).onUploadComplete(({ file }) => {
    console.log("Upload completed:", file);
    // The response to the client is handled by Uploadthing internally,
    // but we can log the file details here
    console.log("File URL:", file.ufsUrl);
    // Note: Returning data here doesn't directly affect the HTTP response
  }),
}
