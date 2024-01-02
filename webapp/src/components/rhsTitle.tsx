// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {Portal} from '@mantine/core';
import React from 'react';
import styled from 'styled-components';

const placeholderID = 'customer-rhs-title-placeholder';

const PlaceholderContainer = styled.div`
    width: 100%;
    height: 100%;
`;

const RHSTitlePlaceholder = () => {
    return (
        <PlaceholderContainer
            id={placeholderID}
        />
    );
};

const RHSTitleRemoteRender = (props: {children: React.ReactNode}) => {
    return (
        <Portal target={'#' + placeholderID}>
            {props.children}
        </Portal>
    );
};

export {
    RHSTitleRemoteRender as RHSTitle,
    RHSTitlePlaceholder,
};
