import { useState } from "react";
import FlashCard from "../ui/flashcard/flashcard";
import PlusButton from "../ui/plusbutton/plusbutton";
import "./main.css";

const Main = () => {
    //state to show or hide flashcard
    const [showFlashCard,setshowFlashCard] = useState<boolean>(false);
    return (
        <>
        {/*Component to add flash card*/}
        <PlusButton titleText="Add flash card" onClick={()=>setshowFlashCard(true)}/>
            <button onClick={()=>location.href = "/saved/falshcards"}>Your saved flashcards</button>
        {showFlashCard && <FlashCard close={()=>setshowFlashCard(false)}/>}
        </>
    )
}
export default Main