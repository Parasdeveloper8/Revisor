import LandingPage from "./landingPage"
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
const App = () =>{
  return (
    <>
    <Router>
      <Routes>
        <Route path="/" element={<LandingPage/>}/>
      </Routes>
    </Router>
    </>
  )
}

export default App
