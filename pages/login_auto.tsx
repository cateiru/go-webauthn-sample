import dynamic from "next/dynamic";
const LoginAutoComlete = dynamic(
  () =>
    import("@/components/LoginAutoComplete").then(
      (modules) => modules.LoginAutoComplete
    ),
  { ssr: false }
);

const AutoLogin = () => {
  return (
    <>
      <LoginAutoComlete />
    </>
  );
};

export default AutoLogin;
