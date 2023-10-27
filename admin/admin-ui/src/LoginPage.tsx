import { useState, useEffect } from 'react';
import './styles/LoginPage.scss'

interface Field {
  name: string;
  type: string;
}

interface UserSchema {
  User: User;
}

interface User {
  fields: Field[];
}




const LoginPage = () => {
  const [userSchema, setUserSchema] = useState<UserSchema>();
  interface FormData {
    [key: string]: string;
  }
  const [formData, setFormData] = useState<FormData>({} as FormData);
  

  useEffect(() => {
    fetch('/api/login-schema')
    .then((res) => {
      if (!res.ok) {
        throw new Error('Failed to fetch schema');
    }
    return res.json();
    })
    .then((schema: UserSchema) => {
      setUserSchema(schema);
    })
    .catch((error) => {
      console.error('Error fetching schema:', error);
    });
    fetch('/api/user-initialization-status')
    .then((res) => {
      if (res.status === 500) {
        window.location.href = '/admin-entry/init';
        return JSON.stringify({status: "User not initialized"});
      }
      if (res.status === 200) {
        return JSON.stringify({status: "User initialized"});
      }
    }
    )
  }, []);
  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData ({
      ...formData,
      [name]: value,
    });
  };
  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
  const temp = async () => {
    e.preventDefault();
    try {
      const response = await fetch('/api/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(formData),
      });

      if (response.status === 200) {
        window.location.href = '/admin-ui/dashboard';
      } else {
        console.error('Login failed');
      }
    } catch (error) {
      console.error('Error:', error);
    }
  }
  temp();
  };

  return (
    <>
    <title>Colibase - Login Page</title>
<div className='background'>
            <div className='login-form'>
                <form onSubmit={handleSubmit} className='login-form'>
                    <h1>Login</h1>
                    {userSchema && (
                        <div>
                            {userSchema.User.fields.map((field: Field) => (
                                <div>
                                  {
                                    RegExp ('email', 'i').test(field.name) && (
                                      <input
                                          name={field.name}
                                          id={field.name}
                                          placeholder={field.name}
                                          type='email'
                                          value={formData[field.name] || ''}
                                          onChange={handleChange}
                                          required
                                          aria-required
                                      />
                                    )
                                  }
                                  {
                                    RegExp ('password', 'i').test(field.name) && (
                                      <input
                                          name={field.name}
                                          id={field.name}
                                          placeholder={field.name}
                                          type='password'
                                          value={formData[field.name] || ''}
                                          onChange={handleChange}
                                          required
                                          aria-required
                                      />
                                    )
                                  }
                                  {
                                    !RegExp ('email', 'i').test(field.name) && !RegExp ('password', 'i').test(field.name) && (
                                      <input
                                          name={field.name}
                                          id={field.name}
                                          type='text'
                                          placeholder={field.name}
                                          value={formData[field.name] || ''}
                                          onChange={handleChange}
                                          required
                                          aria-required
                                      />
                                    )
                                  }
                                </div>
                            ))}
                        </div>
                    )}
                    <button type='submit'>Submit</button>
                </form>
            </div>
    </div>
    </>
  )
}

export default LoginPage
