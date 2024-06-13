import { createContext, useContext } from 'react';

export const SessionContext = createContext({ user: null, setUser: () => {} });

export const useSession = () => {
  const context = useContext(SessionContext);
  if (!context) {
    throw new Error('useSession must be used within a Session');
  }

  const { user, setUser } = context;
  return { user, setUser };
};
