import { useEffect, useState } from "react";
import "./savedFlashCard.css";
import Loader from "../loader/loader";
import {useGlobalContext} from "../../context/GlobalContext";
import {generateQuiz} from "../../utils/quizUtils";
import { useNavigate } from "react-router-dom";

// type to hold each flashcard item
type FlashCardItem = {
  Email: string;
  TopicName: string;
  Time: string;
  FormattedTime: string;
  Data: {
    Heading: string;
    Value: string;
  }[];
  Uid : string;
};

// full API response
type FlashCardApiResponse = {
 flashCardData: FlashCardItem[];
};

const SavedFlashCards = () => {
  const [flashCardItems, setFlashCardItems] = useState<FlashCardItem[]>([]);
  const [fetched,setFetched] = useState<boolean>(false);
  const {generatingQuiz,setIsGenerated} = useGlobalContext();
   const navigate = useNavigate();
  // Fetch all saved flashcards' data
  useEffect(() => {
    fetch("http://localhost:8080/flashcard/get/data", {
      method: "GET",
      credentials: "include",
    })
      .then((res) => res.json())
      .then((data: FlashCardApiResponse) => {
        console.log("Fetched flashcards:", data);
        setFlashCardItems(data.flashCardData); // ✅ correct key
        setFetched(true);
      })
      .catch((err) => {
        console.error("Failed to fetch flashcards:", err);
        setFetched(true);
      });
  }, []);

  return (
  <div className="flashcards-container">
    <h2>Saved Flashcards</h2>
    <button onClick={() => (location.href = "/")}>Home</button>
    {generatingQuiz && <Loader/> }
    {!fetched ? (
      <Loader/>
    ) : !Array.isArray(flashCardItems) ? (
      <p>Data format error.</p>
    ) : flashCardItems.length === 0 ? (
      <p>No flashcards found.</p>
    ) : (
      flashCardItems.map((item, index) => (
        <div key={index} className="flashcard-item">
          <h3>{item.TopicName}</h3>
          <p>{item.FormattedTime}</p>
          <div className="flashcard-data">
          {item.Data.map((card, i) => (
          <div key={i} className="flashcard-row">
          <div className="flashcard-heading">{card.Heading}</div>
          <div className="flashcard-value">{card.Value}</div>
        </div>
  ))}
      </div>
        <button
      className="generate-btn"
      key={index}
      onClick={()=>generateQuiz(item.TopicName,item.Data,item.Uid,navigate,setIsGenerated)}
    >
      Generate Quiz
    </button>
          <hr />
        </div>
      ))
    )}
  </div>
);

};

export default SavedFlashCards;
