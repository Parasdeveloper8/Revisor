import LandingPage from "./landingPage"
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import SavedFlashCards from "./ui/savedFlashCards/savedFlashCard";
import QuizPage from "./ui/quiz/quiz";
const App = () =>{
  return (
    <>
    <Router>
      <Routes>
        <Route path="/" element={<LandingPage/>}/>
        <Route path="/saved/flashcards" element = {<SavedFlashCards/>} />
        <Route path="/quiz" element = {<QuizPage/>}/>
      </Routes>
    </Router>
    </>
  )
}

export default App
