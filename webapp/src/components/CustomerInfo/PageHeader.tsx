import {Title} from '@mantine/core';
import React from 'react';

type Params = {
    text: string;
}
const PageHeader = (
    {text}: Params,
) => {
    return (
        <Title
            order={3}
        >
            {text}
        </Title>
    );
};

export {
    PageHeader,
};
