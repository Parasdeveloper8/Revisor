import { createContext, useContext, useState, type ReactNode } from 'react';

// Define only the global state shape
interface GlobalState {
  email: string;
  setEmail: React.Dispatch<React.SetStateAction<string>>;
}

// Create the context
const GlobalContext = createContext<GlobalState | undefined>(undefined);

// Provider Props
interface GlobalProviderProps {
  children: ReactNode;
}

// Create the provider component
export const GlobalProvider: React.FC<GlobalProviderProps> = ({ children }) => {
  const [email, setEmail] = useState(() => {
    // Load email from localStorage when component mounts
    if (typeof window !== 'undefined') {
      return localStorage.getItem('email') || '';
    }
    return '';
  });

  return (
    <GlobalContext.Provider value={{ email, setEmail }}>
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
