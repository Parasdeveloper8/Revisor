import React, { useState } from 'react';
import './FlashCard.css';

type FlashCardProps = {
    close: React.MouseEventHandler<HTMLButtonElement>;
};

const FlashCard: React.FC<FlashCardProps> = ({ close }) => {
    const [numberOfFields, setNumberOfFields] = useState<number[]>([1]);

    const addField = () => {
        setNumberOfFields((prev) => [...prev, prev.length + 1]);
    };

    const removeField = () => {
        if (numberOfFields.length > 1) {
            setNumberOfFields((prev) => prev.slice(0, -1));
        }
    };

    return (
        <div className="flashcard">
            <button onClick={close} className="close-btn" aria-label="Close">✕</button>
            <div className="flashcard-content">
                <h2 className="title">Create Flashcard</h2>
                <label>
                    <span>Topic Name</span>
                    <input type="text" placeholder='Ex:= Nationalism in India'/>
                </label>
            
                {numberOfFields.map((_, index) => (
                    <div className="field" key={index}>
                        <label>
                            <span>Heading</span>
                            <input type="text" placeholder="Enter heading" />
                        </label>
                        <label>
                            <span>Text</span>
                            <textarea placeholder="Enter text here..." />
                        </label>
                    </div>
                ))}

                <div className="button-group">
                    {numberOfFields.length > 1 && (
                        <button className="secondary-btn" onClick={removeField}>
                            - Remove Field
                        </button>
                    )}
                    <button className="primary-btn" onClick={addField}>
                        + Add Field
                    </button>
                    <button className="primary-btn">
                       Create flashcard
                    </button>
                </div>
            </div>
        </div>
    );
};

export default FlashCard;
