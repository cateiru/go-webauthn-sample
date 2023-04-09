import { Button, Center, useToast } from "@chakra-ui/react";
import {
  parseRequestOptionsFromJSON,
  get,
} from "@github/webauthn-json/browser-ponyfill";
import React from "react";

export const Login = () => {
  const [challenge, setChallenge] =
    React.useState<CredentialRequestOptions | null>(null);
  const toast = useToast();

  React.useEffect(() => {
    if (!challenge) {
      fetch("http://localhost:1323/begin_login", {
        method: "GET",
        mode: "cors",
        credentials: "include",
      })
        .then((res) => res.json())
        .then((data) => {
          const credentialRequest = parseRequestOptionsFromJSON(data);
          setChallenge(credentialRequest);
        })
        .catch((err) => {
          if (err instanceof Error) {
            toast({
              title: "Error",
              description: err.message,
              status: "error",
            });
          }
        });
    }
  }, []);

  const handleLogin = async () => {
    if (!challenge) return;

    let credential: Credential;
    try {
      credential = await get(challenge);
    } catch (e) {
      if (e instanceof Error) {
        toast({
          title: "Error",
          description: e.message,
          status: "error",
        });
      }
      return;
    }

    const credentialJSON = JSON.stringify(credential);

    try {
      const res = await fetch("http://localhost:1323/login", {
        method: "POST",
        mode: "cors",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
        },
        body: credentialJSON,
      });

      if (!res.ok) {
        throw new Error(await res.text());
      }

      const user = await res.json();

      toast({
        title: "Success",
        description: `Hello ${user.name}`,
        status: "success",
      });
    } catch (err) {
      if (err instanceof Error) {
        toast({
          title: "Error",
          description: err.message,
          status: "error",
        });
      }
      return;
    }
  };

  return (
    <Center h="100vh">
      <Button
        onClick={() => {
          handleLogin();
        }}
      >
        Login with WebAuthn
      </Button>
    </Center>
  );
};
