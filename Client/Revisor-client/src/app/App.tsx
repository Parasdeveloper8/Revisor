import LandingPage from "./landingPage"
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import SavedFlashCards from "./ui/savedFlashCards/savedFlashCard";
const App = () =>{
  return (
    <>
    <Router>
      <Routes>
        <Route path="/" element={<LandingPage/>}/>
        <Route path="/saved/flashcards" element = {<SavedFlashCards/>} />
      </Routes>
    </Router>
    </>
  )
}

export default App
