import { useState } from "react";
import FlashCard from "../ui/flashcard/flashcard";
import PlusButton from "../ui/plusbutton/plusbutton";
import "./main.css";
import SavedFlashCards from "../ui/savedFlashCards/savedFlashCard";

const Main = () => {
    //state to show or hide flashcard
    const [showFlashCard,setshowFlashCard] = useState<boolean>(false);
    //state to show saved flashcards
    const [showSavedFlashCard,setShowSavedFlashCard] = useState<boolean>(false);
    return (
        <>
        {/*Component to add flash card*/}
        <PlusButton titleText="Add flash card" onClick={()=>setshowFlashCard(true)}/>
            <button onClick={()=>setShowSavedFlashCard(true)}>Show saved flashcards</button>
            {/*Component to show saved flash card*/}
            {showSavedFlashCard && <SavedFlashCards/> }
        {showFlashCard && <FlashCard close={()=>setshowFlashCard(false)}/>}
        </>
    )
}
export default Main