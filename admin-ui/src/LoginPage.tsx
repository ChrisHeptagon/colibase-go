import { createSignal, onMount } from 'solid-js';
import './styles/LoginPage.css'

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
  const [userSchema, setUserSchema] = createSignal<UserSchema | null>(null);
  interface FormData {
    [key: string]: string;
  }
  let formData: FormData = {};
  

  onMount(() => {
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
  });

  const handleChange = (e: any) => {
    const { name, value } = e.target;
    formData = ({
      ...formData,
      [name]: value,
    });
  };
  const handleSubmit = async (e: any) => {
    e.preventDefault();
    try {
      // Send a POST request with formData to your server
      const response = await fetch('/api/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(formData),
      });

      if (response.status === 200) {
        // Redirect or perform actions for a successful login
      } else {
        // Handle login failure
        console.error('Login failed');
      }
    } catch (error) {
      console.error('Error:', error);
    }
  };

  return (
    <>
    <title>Colibase - Login Page</title>
<div class='background'>
            <div class='login-form'>
                <form onSubmit={handleSubmit} class='login-form'>
                    <h1>Login</h1>
                    {userSchema() && (
                        <div>
                            {userSchema()?.User.fields.map((field) => (
                                <div>
                                  {field.name !== 'password' && field.name !== 'Password' && (
                                    <input
                                        name={field.name}
                                        id={field.name}
                                        type='text'
                                        placeholder={field.name}
                                        value={formData[field.name] || ''}
                                        onChange={handleChange}
                                    />
                                  )}
                                  {
                                    field.name === 'password' || field.name === 'Password' && (
                                      <input
                                          name={field.name}
                                          id={field.name}
                                          placeholder={field.name}
                                          type='password'
                                          value={formData[field.name] || ''}
                                          onChange={handleChange}
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
