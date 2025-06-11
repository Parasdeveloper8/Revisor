import PlusButton from "../ui/plusbutton/plusbutton";
import "./main.css";

const Main = () => {
    return (
        <>
        <PlusButton titleText="Add flash card" onClick={()=>alert("Flash..")}/>
        </>
    )
}
export default Main