 import {type NavigateFunction } from "react-router-dom";

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

  //This function sends data which has to be converted into quiz to backend
export const generateQuiz = (
  topicName:FlashCardItem["TopicName"], data:FlashCardItem["Data"], noteId:FlashCardItem["Uid"],
  navigate:NavigateFunction, setIsGenerated: React.Dispatch<React.SetStateAction<boolean>>,
       )=>{
      //type to hold /generate/quiz response data
      type QuizQuesData = {
         response: {
          QuizId: string;
            Quesopts: {
            Question: string;
            Options: string[];
           }[];
            };
         topic: string;
       }
      
      setIsGenerated(true);
     const api:string = "http://localhost:8080/generate/quiz";
     fetch(api,{
      method:'POST',
      credentials:'include',
      headers:{ 'Content-Type': 'application/json' },
      body:JSON.stringify({"topicName":topicName,"data":data,"noteId":noteId}),
     })
     .then(response=>response.json())
     .then((data :QuizQuesData) => {
          console.log(data);
          setIsGenerated(false);
          //move user to separate quiz page
          navigate('/quiz',{state:{quizData:data}});
     })
     .catch((error) => {
            console.error("Failed to make post request to /generate/quiz: ",error);
            setIsGenerated(false);
         })
}