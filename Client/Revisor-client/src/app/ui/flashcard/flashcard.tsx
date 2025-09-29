import React, { useState } from 'react';
import './FlashCard.css';
import SuccessTick from '../successTick/succTick';
import SpeechRecognition, { useSpeechRecognition } from 'react-speech-recognition';

type FlashCardProps = {
  close: React.MouseEventHandler<HTMLButtonElement>;
};

type FlashCardData = {
  heading: string;
  value: string;
};

type FlashMsg = {
  type: 'error' | 'success';
  message: string;
};

const FlashCard: React.FC<FlashCardProps> = ({ close }) => {
  const [flashCardData, setFlashCardData] = useState<FlashCardData[]>([{ heading: '', value: '' }]);
  const [topic, setTopic] = useState<string>('');
  const [flashMsg, setFlashMsg] = useState<FlashMsg | null>(null);
  const textRefs = React.useRef<(HTMLTextAreaElement | null)[]>([]);

  // This state keeps track of which textarea is currently recording
  const [activeIndex, setActiveIndex] = useState<number | null>(null);

  // Use react-speech-recognition hook
  const { transcript, listening, resetTranscript, browserSupportsSpeechRecognition } = useSpeechRecognition();

  const addField = () => setFlashCardData(prev => [...prev, { heading: '', value: '' }]);
  const removeField = () => flashCardData.length > 1 && setFlashCardData(prev => prev.slice(0, -1));

  const handleChange = (index: number, field: 'heading' | 'value', value: string) => {
    setFlashCardData(prev =>
      prev.map((item, i) => (i === index ? { ...item, [field]: value } : item))
    );
  };

  const sendFlashCardData = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const response = await fetch('http://localhost:8080/flashcard/store/data', {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ topic, flashdata: flashCardData }),
      });
      if (response.ok) {
        setFlashMsg({ type: 'success', message: 'Flashcard created' });
      } else {
        const failMsg = await response.json();
        setFlashMsg({ type: 'error', message: failMsg.info || failMsg.message || 'Something went wrong' });
      }
    } catch (err) {
      console.error('Failed to send data:', err);
    }
  };

  const autoResize = (index:number) => {
    const textArea = textRefs.current[index];
    if (!textArea) return ;
    textArea.style.height = 'auto';
    textArea.style.height = textArea.scrollHeight + 'px';
  };

  // Start or stop speech recognition for a specific textarea
  const toggleListening = (index: number) => {
    if (!browserSupportsSpeechRecognition) {
      alert('Speech recognition not supported in this browser');
      return;
    }
    if (activeIndex === index) {
      SpeechRecognition.stopListening();
      setActiveIndex(null);
      resetTranscript();
    } else {
      resetTranscript();
      setActiveIndex(index);
      SpeechRecognition.startListening({ continuous: true });
    }
    
  };

  // Update textarea in real-time while recording
  React.useEffect(() => {
    if (activeIndex !== null) {
      handleChange(activeIndex, 'value', transcript);
      autoResize(activeIndex);
    }
  }, [transcript, activeIndex]);

  return (
    <div className="flashcard">
      <button onClick={close} className="close-btn" aria-label="Close">✕</button>

      {flashMsg ? (
        <SuccessTick message={flashMsg.message} color={flashMsg.type === 'error' ? 'red' : 'green'} />
      ) : (
        <div className="flashcard-content">
          <h2 className="title">Create Flashcard</h2>
          <form onSubmit={sendFlashCardData}>
            <label>
              <span>Topic Name</span>
              <input
                type="text"
                placeholder="Ex:= Nationalism in India"
                value={topic}
                onChange={e => setTopic(e.target.value)}
              />
            </label>

            {flashCardData.map((_, index) => (
              <div className="field" key={index}>
                <label>
                  <span>Heading</span>
                  <input
                    type="text"
                    placeholder="Enter heading"
                    value={flashCardData[index]?.heading || ''}
                    onChange={e => handleChange(index, 'heading', e.target.value)}
                  />
                </label>
                <label>
                  <span>Text</span>
                  <textarea
                  ref={el => {textRefs.current[index] = el;}}
                    placeholder="Enter text here..."
                    value={flashCardData[index]?.value || ''}
                    onChange={e => {
                      handleChange(index, 'value', e.target.value);
                      autoResize(index);
                    }}
                  />
                </label>
                <button type="button" onClick={() => toggleListening(index)}>
                  {activeIndex === index && listening ? '🛑 Stop' : '🎤 Start'}
                </button>
              </div>
            ))}

            <button className="primary-btn" type="submit">Create flashcard</button>
          </form>

          <div className="button-group">
            {flashCardData.length > 1 && (
              <button className="secondary-btn" onClick={removeField}>- Remove Field</button>
            )}
            <button className="primary-btn" onClick={addField}>+ Add Field</button>
          </div>
        </div>
      )}
    </div>
  );
};

export default FlashCard;
