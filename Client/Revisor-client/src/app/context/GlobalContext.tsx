import React, { createContext, useContext, useState, type ReactNode } from 'react';

type notification = {
   message : string;
   anyInfo : 'success' | 'error' | 'noInfo';
}

type generatingQuiz = boolean;
interface GlobalState {
  name: string;
  email: string;
  token: string;
  tokenExpiry: Date;
  info : notification;
  generatingQuiz:boolean;
  setIsGenerated:React.Dispatch<React.SetStateAction<generatingQuiz>>;
  setInfo : React.Dispatch<React.SetStateAction<notification>>;
  setUserData: React.Dispatch<React.SetStateAction<UserData>>;
}

interface UserData {
  name: string;
  email: string;
  token: string;
  tokenExpiry: Date;
}

// Create the context
const GlobalContext = createContext<GlobalState | undefined>(undefined);

// Provider Props
interface GlobalProviderProps {
  children: ReactNode;
}

// Create the provider component
export const GlobalProvider: React.FC<GlobalProviderProps> = ({ children }) => {
  const [userData, setUserData] = useState<UserData>({
    name: '',
    email: '',
    token: '',
    tokenExpiry: new Date(0),
  });
  const [info,setInfo] = useState<notification>({
     message:'',
     anyInfo:'noInfo'
  });
  const [generatingQuiz,setIsGenerated] = useState<generatingQuiz>(false);
  return (
    <GlobalContext.Provider value={{ ...userData, setUserData,info,setInfo,generatingQuiz,setIsGenerated}}>
      {children}
    </GlobalContext.Provider>
  );
};

// Custom hook to use global context
export const useGlobalContext = () => {
  const context = useContext(GlobalContext);
  if (!context) throw new Error('useGlobalContext must be used inside GlobalProvider');
  return context;
};
