#version 460 core

layout (location = 0) out vec4 FragColor;

layout (location = 2) in FragBlock
{
    vec2 TexCoord;
    vec3 Normal;
    vec3 FragPos;
};

layout (binding = 0) uniform UBvec3s
{
    vec3 lightPos;
    vec3 lightColor;
};

layout (binding = 0) uniform sampler2D texture1;

void main() {
    vec4 ambientLight = vec4(0.2, 0.05, 0.3, 1.0);
    FragColor = ambientLight * texture(texture1, TexCoord);
};
