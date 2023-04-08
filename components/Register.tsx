import React from "react";
import { Box, Button, Center, Input, useToast } from "@chakra-ui/react";
import {
  create,
  parseCreationOptionsFromJSON,
} from "@github/webauthn-json/browser-ponyfill";

const Register = () => {
  const toast = useToast();
  const [name, setName] = React.useState("");

  const handleRegister = async () => {
    let response: Response;
    try {
      response = await fetch(
        `http://localhost:1323/begin_create?name=${name}`,
        {
          method: "GET",
          mode: "cors",
          credentials: "include",
        }
      );
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

    const data = await response.json();

    const credentialData = parseCreationOptionsFromJSON(data);

    let credential: Credential | null;
    try {
      credential = await create(credentialData);
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
    if (!credential) {
      toast({
        title: "Error",
        description: "Failed to create credential",
        status: "error",
      });
      return;
    }

    const credentialJSON = JSON.stringify(credential);
    const createResponse = await fetch("http://localhost:1323/create", {
      method: "POST",
      mode: "cors",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: credentialJSON,
    });

    if (createResponse.status !== 200) {
      toast({
        title: "Failed to create credential",
        description: await createResponse.text(),
        status: "error",
      });
      return;
    }
  };

  return (
    <Center h="100vh">
      <Box>
        <Input
          placeholder="Username"
          defaultValue={name}
          onChange={(v) => setName(v.target.value)}
        />
        <Button
          mt=".5rem"
          onClick={() => {
            handleRegister();
          }}
        >
          Registered WebAuthn
        </Button>
      </Box>
    </Center>
  );
};

export default Register;
