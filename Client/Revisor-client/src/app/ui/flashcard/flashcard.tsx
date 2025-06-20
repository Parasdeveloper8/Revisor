import React, { useState } from 'react';
import './FlashCard.css';
import SuccessTick from '../success-tick/succTick';

type FlashCardProps = {
    close: React.MouseEventHandler<HTMLButtonElement>;
};
type FlashCardData = {
    heading:string;
    value:string;
}

const FlashCard: React.FC<FlashCardProps> = ({ close }) => {
   // const [numberOfFields, setNumberOfFields] = useState<number[]>([1]);
    const [flashCardData,setFlashCardData] = useState<FlashCardData[]>([{heading:'',value:''}]);
    const [topic,setTopic] = useState<string>('');
    const [isCreated,setIsCreated] = useState<boolean>(false); //State to show flashcard is created or not

    const addField = () => {
        //setNumberOfFields((prev) => [...prev, prev.length + 1]);
        //Also create a blank object so that it can be updated later
        setFlashCardData(prev => [...prev,{heading:'',value:''}]);
    };

    const removeField = () => {
        if (flashCardData.length > 1) {
           // setNumberOfFields((prev) => prev.slice(0, -1));
            setFlashCardData(prev =>prev.slice(0,-1));//remove object also
        }
    };
     //this function adds data in flashCardData state
    const handleChange = (index: number, field: 'heading' | 'value', value: string) => {
     setFlashCardData(prev =>
    prev.map((item, i) =>
      i === index ? { ...item, [field]: value } : item
    )
    );}
    console.log(flashCardData);
    return (
        <div className="flashcard">
            <button onClick={close} className="close-btn" aria-label="Close">✕</button>
            {isCreated ? (
                <SuccessTick message='Flashcard created'/>
            ) : (
            <div className="flashcard-content">
                <h2 className="title">Create Flashcard</h2>
                <form onSubmit={(e)=>{e.preventDefault();setIsCreated(true)}}>
                <label>
                    <span>Topic Name</span>
                    <input type="text" placeholder='Ex:= Nationalism in India' value={topic} onChange={e => setTopic(e.target.value)}/>
                </label>
                <br/>
                <br/>
                {flashCardData.map((_, index) => (
                    <div className="field" key={index}>
                        <label>
                            <span>Heading</span>
                            <input
                type="text"
              placeholder="Enter heading"
              value={flashCardData[index]?.heading || ''}
              onChange={(e) => handleChange(index, 'heading', e.target.value)}
                      />
                        </label>
                        <label>
                            <span>Text</span>
                            <textarea
                       placeholder="Enter text here..."
                      value={flashCardData[index]?.value || ''}
                      onChange={(e) => handleChange(index, 'value', e.target.value)}
                           />
                        </label>
                    </div>
                ))}
                <br/>
                <button className="primary-btn" type='submit'>
                       Create flashcard
                    </button>
                </form>
                <div className="button-group">
                    {flashCardData.length > 1 && (
                        <button className="secondary-btn" onClick={removeField}>
                            - Remove Field
                        </button>
                    )}
                    <button className="primary-btn" onClick={addField}>
                        + Add Field
                    </button>
                </div>
            </div>
            )}
        </div>
    );
};

export default FlashCard;
