import { useEffect, useState } from 'react';
import { SessionContext } from './hooks/useSession';
import { baseUrl } from './constants/api';

export default function Session({ children }) {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const getUser = async () => {
      try {
        const res = await fetch(`${baseUrl}/user`, {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json'
          }
        });

        if (res.status !== 200) {
          throw new Error('Failed to get user');
        }

        const data = await res.json();
        setUser(data);
      } catch (error) {
        console.log(error);
      } finally {
        setLoading(false);
      }
    };

    getUser();
  }, []);

  const handleSetUser = newUser => {
    setUser(newUser);
  };

  if (loading) {
    return <div>Loading...</div>;
  }

  return <SessionContext.Provider value={{ user, setUser: handleSetUser }}>{children}</SessionContext.Provider>;
}
