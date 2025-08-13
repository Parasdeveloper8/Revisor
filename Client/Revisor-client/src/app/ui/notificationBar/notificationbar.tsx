// src/components/NotificationBar.tsx
import React, { useEffect, useState } from "react";
import "./notificationBar.css";
import { useGlobalContext } from "../../context/GlobalContext";

interface NotificationBarProps {
  message: string;
  type?: "success" | "error";
  duration?: number; // in ms
}

const NotificationBar: React.FC<NotificationBarProps> = ({
  message,
  type = "success",
  duration = 3000,
}) => {
  const [visible, setVisible] = useState(true);
  const {setInfo} = useGlobalContext();
  useEffect(() => {
    const timer = setTimeout(() =>{
      setVisible(false);
      setInfo({
        message : '',
        anyInfo : 'noInfo'
      })
    }, duration);
    return () => clearTimeout(timer);
  }, [duration,setInfo]);

  return visible ? (
    <div className={`notification-bar ${type}`}>
      {message}
    </div>
  ) : null;
};

export default NotificationBar;
