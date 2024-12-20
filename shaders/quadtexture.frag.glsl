#version 460 core

layout (location = 0) out vec4 FragColor;

layout (location = 1) in InBlock
{
    vec2 TexCoord;
};

layout (binding = 1) uniform sampler2D texture1;

void main() {
    FragColor = texture(texture1, TexCoord);
    // FragColor = vec4(tex.x, tex.y, 1.0, 1.0);
};
