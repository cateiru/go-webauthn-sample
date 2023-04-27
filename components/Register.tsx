import React from "react";
import {
  Box,
  Button,
  Center,
  Heading,
  Input,
  Link,
  ListItem,
  UnorderedList,
  useToast,
} from "@chakra-ui/react";
import {
  create,
  parseCreationOptionsFromJSON,
} from "@github/webauthn-json/browser-ponyfill";
import NextLink from "next/link";

export const Register = () => {
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
        <Heading>WebAuthn を登録</Heading>
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
        <UnorderedList mt="1rem">
          <ListItem>
            <Link href="/login" as={NextLink}>
              ボタンでログイン
            </Link>
          </ListItem>
          <ListItem>
            <Link href="/login_auto" as={NextLink}>
              Conditional UIでログイン
            </Link>
          </ListItem>
        </UnorderedList>
      </Box>
    </Center>
  );
};
