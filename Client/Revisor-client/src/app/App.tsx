import LandingPage from "./landingPage";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import SavedFlashCards from "./ui/savedFlashCards/savedFlashCard";
import QuizPage from "./ui/quiz/quiz";
import ResultPage from "./ui/result/result";
const App = () =>{
  return (
    <>
    <Router>
      <Routes>
        <Route path="/" element={<LandingPage/>}/>
        <Route path="/saved/flashcards" element = {<SavedFlashCards/>} />
        <Route path="/quiz" element = {<QuizPage/>}/>
        <Route path="/result" element={<ResultPage/>}/>
      </Routes>
    </Router>
    </>
  )
}

export default App
