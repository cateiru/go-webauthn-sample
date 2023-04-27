import {
  Button,
  Center,
  Input,
  useToast,
  Divider,
  Box,
} from "@chakra-ui/react";
import {
  parseRequestOptionsFromJSON,
  get,
  supported,
} from "@github/webauthn-json/browser-ponyfill";
import React from "react";
import { browserName } from "react-device-detect";

export const LoginAutoComplete = () => {
  const toast = useToast();
  const [run, setRun] = React.useState(false);

  const isFirefox = browserName === "Firefox";

  React.useEffect(() => {
    if (run) return;

    const abort = new AbortController();

    const f = async () => {
      if (!supported()) {
        return;
      }

      // FirefoxはConditional UIに対応していないので、一旦無視
      if (isFirefox) {
        return;
      }

      console.log("setup");

      const challenge = await fetch("http://localhost:1323/begin_login", {
        method: "GET",
        mode: "cors",
        credentials: "include",
      });
      let data = parseRequestOptionsFromJSON(await challenge.json());

      // See also: https://github.com/w3c/webauthn/wiki/Explainer:-WebAuthn-Conditional-UI
      // W3C上ではConditional UIはdraftなので一旦anyで型を無視している
      // 実装は、メジャーなブラウザ（Chrome, Safari, FireFox）で対応しているっぽい
      (data as any).mediation = "conditional";

      data.signal = abort.signal;

      let credential: Credential;
      try {
        credential = await get(data);
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
    setRun(false);
    f();

    return () => {
      // AbortControllerを使って、WebAuthnの処理を中断する
      abort.abort();
    };
  }, [run]);

  // FireFoxのためのやつ
  // ボタンを押したときにWebAuthnのログインを実行する
  const handleLogin = async () => {
    if (!supported()) {
      return;
    }

    const challenge = await fetch("http://localhost:1323/begin_login", {
      method: "GET",
      mode: "cors",
      credentials: "include",
    });
    let data = parseRequestOptionsFromJSON(await challenge.json());

    let credential: Credential;
    try {
      credential = await get(data);
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
      <Box>
        <Input autoComplete="username webauthn" w="300px" />
        {isFirefox && (
          <Button
            onClick={() => {
              handleLogin();
            }}
          >
            Login With WebAuthn
          </Button>
        )}
      </Box>
    </Center>
  );
};
