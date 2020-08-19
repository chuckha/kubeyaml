import React from 'react';
import Tabs from "react-bulma-components/lib/components/tabs";

class Error extends React.Component {
    render() {
        return (
            <Tabs.Tab active={this.props.active} onClick={this.props.clickHandler}>{this.props.version} Errors</Tabs.Tab>
        )
    }
}

export default Error
