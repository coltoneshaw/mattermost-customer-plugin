import React from 'react';

type Params = {
    message: string;
}

/**
 * @description Flex centered full height and width component with a message
 */
const CenteredText = (
    {
        message,
    }: Params,
) => {
    return (
        <div
            style={{
                display: 'flex',
                justifyContent: 'center',
                alignItems: 'center',
                width: '100%',
                height: '100%',
            }}
        >
            {message}
        </div>
    );
};

export {
    CenteredText,
};
