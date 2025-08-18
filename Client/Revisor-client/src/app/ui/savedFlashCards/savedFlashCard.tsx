import { useEffect, useState } from "react";
import "./savedFlashCard.css";
import Loader from "../loader/loader";
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


//type to hold /generate/quiz response data
type QuizQuesData = {
   Choices : {
         Message : {
            Content : string;
         };
      }[];
    }

type BackendResWrapper = {
  response: QuizQuesData;
  topic: string;
};

const SavedFlashCards = () => {
  const [flashCardItems, setFlashCardItems] = useState<FlashCardItem[]>([]);
  const [fetched,setFetched] = useState<boolean>(false);
  const [generatingQuiz,setIsGenerated] = useState<boolean>(false);
   
  //This function sends data which has to be converted into quiz to backend
const generateQuiz = (topicName:FlashCardItem["TopicName"],data:FlashCardItem["Data"])=>{
      setIsGenerated(true);
     const api:string = "http://localhost:8080/generate/quiz";
     fetch(api,{
      method:'POST',
      credentials:'include',
      headers:{ 'Content-Type': 'application/json' },
      body:JSON.stringify({"topicName":topicName,"data":data}),
     })
     .then(response=>response.json())
     .then((data :BackendResWrapper) => {
          console.log(data);
          setIsGenerated(false);
          //move user to separate quiz page
          location.href = `/quiz?questions=${data.response.Choices[0].Message.Content}`;
     })
     .catch((error) => {
            console.error("Failed to make post request to /generate/quiz: ",error);
            setIsGenerated(false);
         })
}

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
      onClick={()=>generateQuiz(item.TopicName,item.Data)}
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
