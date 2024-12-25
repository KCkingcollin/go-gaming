#version 460 core

layout (location = 0) in vec3 aPos;
layout (location = 1) in vec2 aTexCoord;
layout (location = 2) in vec3 aNormal;

layout (location = 3) out FragBlock
{
    vec3 Normal;
    vec3 FragPos;
    vec2 TexCoord;
};

layout (binding = 0) uniform UBMat4s
{
    mat4 model;
    mat4 view;
    mat4 projection;
};

void main()
{
    Normal = aNormal;
    FragPos = vec3(model * vec4(aPos, 1.0));
    gl_Position = projection * view * vec4(FragPos, 1.0);
    TexCoord = vec2(aTexCoord.x, 1.0 - aTexCoord.y);
};
