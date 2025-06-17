import { useState } from "react";
import FlashCard from "../ui/flashcard/flashcard";
import PlusButton from "../ui/plusbutton/plusbutton";
import "./main.css";

const Main = () => {
    //state to show or hide flashcard
    const [showFlashCard,setshowFlashCard] = useState<boolean>(false);
    return (
        <>
        <PlusButton titleText="Add flash card" onClick={()=>setshowFlashCard(true)}/>
        {showFlashCard && <FlashCard close={()=>setshowFlashCard(false)}/>}
        </>
    )
}
export default Main