import { useSearchParams } from "react-router-dom";
import "./quiz.css";

const QuizPage = () => {
  const [searchParam] = useSearchParams();
  const questions = searchParam.get("questions");
  const decoded = questions ? decodeURIComponent(questions) : "";
  const marks: number = 2;

  const decodedQues = decoded
    .split(/\d+\.\s+/) // Split by 1. , 2. , etc.
    .filter(q => q.trim() !== "");

  return (
    <>
      <h1>Quiz Page</h1>
      <button id="homeBtn" onClick={() => (location.href = "/")}>Home</button>
      <ul>
        {decodedQues.map((q, index) => (
          <li key={index} className="question-item">
            <div className="question-text">
              {index + 1}. {q} <b>[{marks}]</b>
            </div>
            <textarea
              className="question-input"
              placeholder="Your answer..."
              rows={1}
              onInput={(e) => {
                const target = e.target as HTMLTextAreaElement;
                target.style.height = "auto";
                target.style.height = target.scrollHeight + "px";
              }}
            ></textarea>
          </li>
        ))}
      </ul>
    </>
  );
};

export default QuizPage;
