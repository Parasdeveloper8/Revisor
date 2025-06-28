import React from 'react';
import './successTick.css';

type Props = {
    message?: string;                                                                           
    color:string;
};

const SuccessTick: React.FC<Props> = ({ message,color}) => {
    return (
        <div className="success-container" style={{"color":color}}>
            {color === "red" ? (<div className="cut">&times;</div> ) : (<div className="tick">&#10004;</div> )}
            <div className="message">{message}</div>
        </div>
    );
};

export default SuccessTick;
