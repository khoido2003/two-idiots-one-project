// index.js
import express from "express";
import { createRouteHandler } from "uploadthing/express";
import { uploadRouter } from "./api/uploadthing.js";
import dotenv from "dotenv";
import cors from "cors";

dotenv.config();

const app = express();
const port = 8145;

// Enable CORS
app.use(
  cors({
    origin: "http://localhost:5173", // SvelteKit dev server origin
    methods: ["GET", "POST"],
    allowedHeaders: ["Content-Type"],
  }),
);

app.use(
  "/api/uploadthing",
  createRouteHandler({
    router: uploadRouter,
    config: {
      // Add your Uploadthing token if required
       token: process.env.UPLOADTHING_TOKEN
    },
  }),
);

app.listen(port, () => {
  console.log("Listening on port " + port);
});
