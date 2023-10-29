import { Button, Grid, Stack, TextField } from "@mui/material";
import { useEffect, useMemo, useState } from "react";


interface DefaultField {
  form_type: string;
  name: string;
}

interface FormInfo {
  [key: string]: string;
}

export default function LoginForm({
      headerText,
}: {
    headerText: string;
}) {
  const [formData, setFormData] = useState<FormInfo>({} as FormInfo);
  const [userSchema, setUserSchema] = useState<DefaultField[]>();
  const [isSubmitted, setIsSubmitted] = useState<boolean>(false);
  const tempArray: string[] = useMemo(() => {
    return [];
  }
  , []);
   const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setIsSubmitted(true);
      const temp = async () => {
      if (window.location.pathname === "/admin-entry/init") {
      try {
        const response = await fetch("/api/init-login", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(formData),
        });
        if (response.status === 200) {
          window.location.href = "/admin-entry/login";
        } else {
          console.error("Login failed");
          response.text().then((text) => {
            if (!tempArray.includes(text)) {
              tempArray.push(`${text}`);
            } else if (tempArray.includes(text)) {
              tempArray.filter((error) => error !== text);
            }
          }
          );
        }
      } catch (error) {
        console.error("Error:", error);
      }
    } else if ( window.location.pathname === "/admin-entry/login") {
      try {
        const response = await fetch("/api/login", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(formData),
        });

        if (response.status === 200) {
          window.location.href = "/admin-ui/dashboard";
        } else {
          console.error("Login failed");
          response.text().then((text) => {
            if (!tempArray.includes(text)) {
              tempArray.push(`${text}`);
            } else if (tempArray.includes(text)) {
              tempArray.filter((error) => error !== text);
            }
          }
          );
        }
      } catch (error) {
        console.error("Error:", error);
      }
    }
  };
    temp();
  
  }
  useEffect(() => {
    fetch("/api/login-schema")
    .then((res) => {
      if (!res.ok) {
        throw new Error("Failed to fetch schema");
      }
      return res.json();
    })
    .then((schema: DefaultField[]) => {
      setUserSchema(schema)
    })
    .catch((error) => {
      console.error("Error fetching schema:", error);
    });
  }
  , []);
  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData({
      ...formData,
      [name]: value,
    });
  };

  return (
    <>
      <Grid className="login-form">
        <form onSubmit={handleSubmit} className="login-form"
        noValidate
        autoComplete="on"
        >
          <h1>{headerText}</h1>
          <Stack spacing={2}
          >
          { userSchema && (
            <div>
              {userSchema.map((field) => (
                <>
                    <TextField
                      name={field.name}
                      id={field.name}
                      type={field.form_type}
                      value={formData[field.name] || ""}
                      onChange={handleChange}
                      variant="outlined"
                      label={field.name}
                      aria-required="true"
                      error={isSubmitted && !formData[field.name]}
                />
                </>
              ))}
            </div>
          )}
          </Stack>
          <Button type="submit"
          variant="contained"
          >Submit</Button>
        </form>
      </Grid>
    </>
  );
}

