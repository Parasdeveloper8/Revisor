import React from "react";

type Props = {
    Title:string;
    Text:string;
}
const WebShareButton: React.FC<Props> = ({Title,Text}) => {
  const handleShare = async () => {
    if (navigator.share) {
      try {
        await navigator.share({
          title: Title,
          text: Text,
          url: window.location.href,
        });
        console.log("Content shared successfully!");
      } catch (error) {
        console.error("Error sharing:", error);
      }
    } else {
      alert("Web Share API is not supported in this browser.");
    }
  };

  return (
    <button
      onClick={handleShare}
      className="px-4 py-2 rounded-lg bg-blue-600 text-white hover:bg-blue-700"
    >
      Share
    </button>
  );
};

export default WebShareButton;
