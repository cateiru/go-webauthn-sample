import {
  Center,
  ListItem,
  UnorderedList,
  Link,
  Heading,
  Box,
} from "@chakra-ui/react";
import NextLink from "next/link";

const Index = () => {
  return (
    <Center h="100vh">
      <Box>
        <Heading>WebAuthnサンプル</Heading>
        <UnorderedList mt="1rem">
          <ListItem>
            <Link href="/register" as={NextLink}>
              新規登録
            </Link>
          </ListItem>
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

export default Index;
