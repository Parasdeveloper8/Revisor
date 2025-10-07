import { useLocation, useNavigate } from "react-router-dom";
import "./quiz.css";
import { useState ,useEffect} from "react";
import { formatTime} from "../../utils/timeUtils";

const QuizPage = () => {
  type QuesOpts = {
    Question: string;
    Options: string[];
  };

  type QuizQuesData = {
    response: {
      QuizId: string;
      Quesopts: QuesOpts[];
    };
    topic: string;
  };
  const location = useLocation();
  const { quizData } = location.state || {};
  const [isOptionMissing,setOptionMissing] = useState<boolean>(false);
  const navigate = useNavigate();
  const [seconds, setSeconds] = useState<number>(0);

   // Start timer when quiz loads
  useEffect(() => {
    const interval = setInterval(() => {
      setSeconds((prev) => prev + 1);
    }, 1000);

    // Cleanup when component unmounts
    return () => clearInterval(interval);
  }, []);

  useEffect(() => {
  (quizData as QuizQuesData).response.Quesopts.forEach((q: QuesOpts) => {
    if (!Array.isArray(q.Options) || q.Options.length === 0) {
      setOptionMissing(true);
    }
  });
   }, [quizData]);

  // Check if quizData and required fields are valid
  const isValidQuizData =
    quizData &&
    quizData.response &&
    Array.isArray(quizData.response.Quesopts) &&
    quizData.response.Quesopts.length > 0;

  const [answers, setAnswers] = useState<(string | null)[]>(
    isValidQuizData
      ? Array(quizData.response.Quesopts.length).fill(null)
      : []
  );

  const handleOptionSelect = (qIndex: number, option: string) => {
    const updatedAnswers = [...answers];
    updatedAnswers[qIndex] = option;
    setAnswers(updatedAnswers);
  };

  if (!isValidQuizData) {
    return (
      <div className="error-boundary">
        <h2>It's not you, it's us.</h2>
        <p>Something went wrong with loading the quiz questions or options.</p>
        <button onClick={() => navigate("/")}>Go Back Home</button>
      </div>
    );
  }

  //send request to server to delete existing quiz belonging to current quizID
  const deleteQuiz = () =>{
      const api = "http://localhost:8080/delete/quiz";
       fetch(api,{
       method:'DELETE',
       credentials:'include',
       headers:{ 'Content-Type': 'application/json' },
       body: JSON.stringify({"quizId":(quizData as QuizQuesData).response.QuizId})
     }).then(response => {
      if(response.ok){
        console.log(`Quiz[${(quizData as QuizQuesData).response.QuizId}] is deleted`);
        setOptionMissing(false);
      }else{
        console.error("Error in deleting quiz");
        setOptionMissing(true);
      }
     })
     .catch(error =>{
            console.error("Failed to make request to  /delete/quiz: ",error);
            setOptionMissing(true);
          })
  }
  if(isOptionMissing) deleteQuiz();
  const evaluateMarks = (quizId:string) =>{
     const api = "http://localhost:8080/evaluate/quiz";
     fetch(api,{
       method:'POST',
       credentials:'include',
       headers:{ 'Content-Type': 'application/json' },
       body: JSON.stringify({"userAnswers":answers,"quizId":quizId,"time":seconds})
     })
       .then(response => response.json())
       .then((data)=>{
        const result = {
          Marks:data.marks,
          TotalMarks:data["total marks"],
          Time:data.time
        };
            navigate('/result',{state:{result}});
       })
         .catch(error =>{
            console.error("Failed to make request to  /evaluate/quiz: ",error);
          })
  }

  return (
    <>
      <h1>Quiz Page</h1>
      <button id="homeBtn" onClick={() => navigate("/")}>
        Home
      </button>

     <div className="timer">
        <h3>{formatTime(seconds)}</h3>
      </div>

      <div className="quiz-container">
        <h3>Topic: {quizData.topic}</h3>

        {(quizData as QuizQuesData).response.Quesopts.map(
          (q: QuesOpts, qIndex: number) => (
            <div key={qIndex} className="question-card">
              <h4>
                Q{qIndex + 1}. {q.Question}
              </h4>
              <div className="options">
                {Array.isArray(q.Options) && q.Options.length > 0 ? (
                  q.Options.map((opt: string, oIndex: number) => (
                    <label key={oIndex} className="option-label">
                      <input
                        type="radio"
                        name={`question-${qIndex}`}
                        value={opt}
                        checked={answers[qIndex] === opt}
                        onChange={() =>
                          handleOptionSelect(qIndex, opt)
                        }
                      />
                      <span>{opt}</span>
                    </label>
                  ))
                ) : (
                  <p>It's not you, it's us. Try it again.</p>
                )
              }
              </div>
            </div>
          )
        )}
      </div>

      <button className="evaluate-btn" onClick={()=>evaluateMarks((quizData as QuizQuesData).response.QuizId)}>Evaluate marks</button>
    </>
  );
};

export default QuizPage;
