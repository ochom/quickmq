import { useRef, useState } from 'react';
import './App.css';
import './Login.css';
import { useSession } from './hooks/useSession';
import { baseUrl } from './constants/api';

function Login() {
  const { setUser } = useSession();
  const usernameRef = useRef(null);
  const passwordRef = useRef(null);
  const [loading, setLoading] = useState(false);

  const submit = async e => {
    e.preventDefault();
    setLoading(true);
    try {
      const res = await fetch(`${baseUrl}/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          username: usernameRef.current.value,
          password: passwordRef.current.value
        })
      });

      if (res.status !== 200) {
        throw new Error('Invalid username or password');
      }

      const data = await res.json();
      setUser(data);
    } catch (error) {
      alert(error?.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <form className='login-box' onSubmit={submit}>
      <div className='lobo-box'>
        <b>QuickMQ</b>
      </div>
      <div className='form-area'>
        <div className='input-area'>
          <label htmlFor='username'>Username</label>
          <input ref={usernameRef} type='text' id='username' name='username' required />
        </div>
        <div className='input-area'>
          <label htmlFor='password'>Password</label>
          <input ref={passwordRef} type='password' id='password' name='password' required />
        </div>
        <div className='button-area'>
          <button id='login' type='submit' disabled={loading}>
            Login
          </button>
        </div>
      </div>
    </form>
  );
}

function App() {
  const { user } = useSession();

  if (!user) return <Login />;

  return (
    <div className='App'>
      <h1>Dashboard</h1>
    </div>
  );
}

export default App;
