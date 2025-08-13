import { useState } from "react";
import FlashCard from "../ui/flashCard/flashcard";
import PlusButton from "../ui/plusButton/plusbutton";
import "./main.css";

const Main = () => {
    //state to show or hide flashcard
    const [showFlashCard,setshowFlashCard] = useState<boolean>(false);
    return (
        <>
        {/*Component to add flash card*/}
        <PlusButton titleText="Add flash card" onClick={()=>setshowFlashCard(true)}/>
        <button onClick={()=>location.href = "/saved/flashcards"} className="ysf-btn">Your saved flashcards</button>
        {showFlashCard && <FlashCard close={()=>setshowFlashCard(false)}/>}
        </>
    )
}
export default Main