import { useLocation, useNavigate } from 'react-router-dom';
import Trophy from '../../../assets/animations/trophy.mp4';
import "./result.css";
const ResultPage = () =>{
     const location = useLocation();
     const navigate = useNavigate();
     const { result } = location.state || {};
      type QuizResult = {
         Marks : number;
         TotalMarks : number;
       }
       const isValidResults = result && typeof (result as QuizResult).Marks === "number" && typeof (result as QuizResult).TotalMarks === "number";
       if(!isValidResults){
          return (
        <div className="error-boundary">
        <h2>It's not you, it's us.</h2>
        <p>Something went wrong with loading the Results</p>
        <button onClick={() => navigate("/")}>Go Back Home</button>
      </div>
    );
       }
    return (
        <>
        <button id="homeBtn" onClick={() => navigate("/")}>
        Home
      </button>
        <div id="resultDiv">
         <video 
                src={Trophy} 
                autoPlay 
                loop 
                muted 
                playsInline
                style={{ width: "300px", height: "auto", display: "block", margin: "0 auto" }}
            />
            <p>Your score is <b>{(result as QuizResult).Marks}</b></p>
            </div>
        </>
    )
}
export default ResultPage