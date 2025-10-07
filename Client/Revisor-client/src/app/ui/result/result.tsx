import { useLocation, useNavigate } from "react-router-dom";
import Trophy from "../../../assets/animations/trophy.mp4";
import "./result.css";
import { useEffect } from "react";
import { formatTime } from "../quiz/quiz";

const ResultPage = () => {
  const location = useLocation();
  const navigate = useNavigate();
  const { result } = location.state || {};

  type QuizResult = {
    Marks: number;
    TotalMarks: number;
    Time:number;
  };

  const isValidResults =
    result &&
    typeof (result as QuizResult).Marks === "number" &&
    typeof (result as QuizResult).TotalMarks === "number";

     useEffect(() => {
    document.body.classList.add("result-body");
    return () => {
      document.body.classList.remove("result-body");
    };
  }, [location.pathname]);

  if (!isValidResults) {
    return (
      <div className="result-error">
        <h2>It's not you, it's us.</h2>
        <p>Something went wrong with loading the Results</p>
        <button onClick={() => navigate("/")}>Go Back Home</button>
      </div>
    );
  }

  return (
    <>
      <button className="result-homeBtn" onClick={() => navigate("/")}>
        Home
      </button>
      <div className="result-container">
        <video
          src={Trophy}
          autoPlay
          loop
          muted
          playsInline
          className="result-trophy"
        />
        <p className="result-score">
          Your score is{" "}
          <b>
            {(result as QuizResult).Marks}/{(result as QuizResult).TotalMarks}
          </b>
        </p>
        {/**<WebShareButton Text="My marks" Title="My marks"/>**/}
        <p>Time taken is{" "}<b>{formatTime((result as QuizResult).Time)}</b></p>
      </div>
    </>
  );
};

export default ResultPage;
