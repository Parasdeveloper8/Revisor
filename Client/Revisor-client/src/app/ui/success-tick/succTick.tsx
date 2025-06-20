import React from 'react';
import './successTick.css';

type Props = {
    message?: string;
};

const SuccessTick: React.FC<Props> = ({ message = "Success!" }) => {
    return (
        <div className="success-container">
            <div className="tick">&#10004;</div>
            <div className="message">{message}</div>
        </div>
    );
};

export default SuccessTick;
