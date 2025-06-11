import React from "react";
import "./plusbutton.css";

type PlusButtonProps = {
   onClick : React.MouseEventHandler<HTMLButtonElement> ;
   titleText : string;
};

const PlusButton: React.FC<PlusButtonProps> = ({onClick ,titleText}) => {
  return (
    <button className="plus-button" onClick={onClick} title= {titleText}>
      +
    </button>
  );
};

export default PlusButton;
