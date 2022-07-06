import { useContext } from "react";
import { AccountInfoContext } from "../infoStore/accountInfoStore";
import { useNavigate } from "react-router-dom";
import { Button } from "@chakra-ui/react";
import { observer } from "mobx-react-lite";

const ClanManageButton = observer(() => {
    const accountInfo = useContext(AccountInfoContext);
    const navigate = useNavigate();
    if (accountInfo.permission != 1) {
        return null;
    }
    return (<Button onClick={() => navigate("/manage", { replace: true })}>
        Manage Clan
    </Button>)
});

export default ClanManageButton;