export {};

declare global {
  interface Navigator {
    share?: (data: {
      title?: string;
      text?: string;
      url?: string;
      files?: File[];
    }) => Promise<void>;

    canShare?: (data?: { files?: File[] }) => boolean;
  }
}
