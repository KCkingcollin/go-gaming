#version 460 core

layout (location = 0) out vec4 FragColor;

layout (location = 3) in FragBlock
{
    vec3 Normal;
    vec3 FragPos;
    vec2 TexCoord;
};

layout (binding = 1) uniform UBVec3s
{
    vec3 lightPos;
    vec3 lightColor;
};

// uniform vec3 lightPos;
// uniform vec3 lightColor;

layout (binding = 0) uniform sampler2D texture1;

void main() {
    //ambient
    vec4 ambientLight = vec4(0.3, 0.3, 0.3, 1.0);
    
    //diffuse
    vec3 lightDir = normalize(lightPos - FragPos);
    float diff = max(dot(Normal, lightDir), 0.0);
    vec3 diffuse = diff * lightColor;

    FragColor = (ambientLight + vec4(diffuse, 1.0)) * texture(texture1, TexCoord);
};
