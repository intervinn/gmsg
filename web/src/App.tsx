import { Route, Switch } from "wouter";
import AuthorizePage from "./pages/AuthorizePage";
import ChatPage from "./pages/ChatPage";


export default function App() {
  return (
    <>
      <Switch>
        <Route path={"/"} component={ChatPage}/>
        <Route path={"/auth"} component={AuthorizePage}/>
      </Switch>
    </>
  )
}