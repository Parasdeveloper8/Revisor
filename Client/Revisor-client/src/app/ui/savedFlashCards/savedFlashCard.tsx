import { useEffect, useState } from "react"
//type to hold data
type FlashCardItem = {
  Email: string;
  TopicName: string;
  Time: string;
  FormattedTime: string;
  Data: {
    Heading: string;
    Value: string;
  }[];
};
const SavedFlashCards = () =>{
  const [flashCardItems , setFlashCardItems] = useState<FlashCardItem[]>([]);
  //Fetch all saved flashcards' data
   useEffect(()=>{
      fetch("http://localhost:8080/flashcard/get/data",{
        method:"GET",
        credentials : "include",
      })
      .then(res => res.json())
      .then((data:FlashCardItem[]) => {
        //set data in state
         setFlashCardItems(data);
      })
   },[]);
     return(
         //render cards
          <div>
      <h2>Saved Flashcards</h2>
      <p>{flashCardItems.length}</p>
      { /*
      {flashCardItems.length === 0 ? (
        <p>No flashcards found.</p>
      ) : (
        flashCardItems.map((item, index) => (
          <div key={index}>
            <h3>{item.TopicName}</h3>
            <p>{item.FormattedTime}</p>
            {item.Data.map((card, i) => (
              <div key={i}>
                <strong>{card.Heading}</strong>: {card.Value}
              </div>
            ))}
            <hr />
          </div>
        ))
      )}
      */}
    </div> 
    )
}
export default SavedFlashCards